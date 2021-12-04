package nimgobus

import (
	"math"
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten/v2/audio"
)

// Nimbus sound synthesis!

const (
	sampleRate   = 44100.0
	baseFreq     = 440.0
	maxAmplitude = 10000
	slew         = 5000
)

type envelope struct {
	attackTime   int
	attackLevel  int
	decayTime    int
	decayLevel   int
	sustainTime  int
	sustainLevel int
	releaseTime  int
}

type noteItem struct {
	note      int
	note2     int
	volume    int
	envelope  envelope
	glissando map[int]int
}

type voice struct {
	audioContext     *audio.Context
	player           *audio.Player
	note             int
	note2            int
	glissando        map[int]int
	selectedEnvelope int
	envelope         envelope
	envelopePosition int
	released         bool
	muQueue          sync.Mutex
	queue            []noteItem
}

// Add a note to the queue if the queue is not full
func (v *voice) Add(item noteItem) bool {
	v.muQueue.Lock()
	if len(v.queue) >= 6 {
		v.muQueue.Unlock()
		return false
	} else {
		v.queue = append(v.queue, item)
		v.muQueue.Unlock()
		return true
	}
}

// Play the oldest note in the queue if no other notes are playing
func (v *voice) Play(n *Nimbus) {
	v.muQueue.Lock()
	if len(v.queue) > 0 {
		if !v.released {
			v.muQueue.Unlock()
			return
		}
		item := v.queue[0]
		v.queue = v.queue[1:]
		v.muQueue.Unlock()
		v.note = item.note
		v.note2 = item.note2
		v.glissando = item.glissando
		v.envelope.attackTime = item.envelope.attackTime
		v.envelope.attackLevel = item.envelope.attackLevel
		v.envelope.decayTime = item.envelope.decayTime
		v.envelope.decayLevel = item.envelope.decayLevel
		v.envelope.sustainTime = item.envelope.sustainTime
		v.envelope.sustainLevel = item.envelope.sustainLevel
		v.envelope.releaseTime = item.envelope.releaseTime
		if item.volume > -1 {
			v.envelope.sustainLevel = item.volume
		}
		v.Trigger()
		for !v.released {
			time.Sleep(10 * time.Millisecond)
		}
	}
}

// Trigger triggers the voice's envelop generator and amplifier
func (v *voice) Trigger() {
	v.released = false
	v.envelopePosition = 0
}

// stream is an infinite stream of 440 Hz sine wave.
type stream struct {
	position  int64
	remaining []byte
	voice     *voice
	printOnce bool
}

// Convert a Note number into its frequency
func frequencyOfNote(note int) float64 {
	distance := note - 157
	freq := baseFreq * math.Pow(2, (float64(distance)/12.0))
	if freq < 1.0 {
		freq = 1.0
	}
	return freq
}

// Read is io.Reader's Read.
//
// Read returns synthesized audio for the voice
func (s *stream) Read(buf []byte) (int, error) {
	if len(s.remaining) > 0 {
		n := copy(buf, s.remaining)
		s.remaining = s.remaining[n:]
		return n, nil
	}

	var origBuf []byte
	if len(buf)%4 > 0 {
		origBuf = buf
		buf = make([]byte, len(origBuf)+4-len(origBuf)%4)
	}

	length := sampleRate / frequencyOfNote(s.voice.note)
	amplitudeScale := maxAmplitude / 16
	p := s.position / 4
	for i := 0; i < len(buf)/4; i++ {
		if !s.printOnce {
			s.printOnce = true
		}
		var b int16
		w := math.Sin(2 * math.Pi * float64(p) / length)
		// Convert sine wave to square
		if w > 0 {
			w = 1
		} else {
			w = -1
		}
		// Apply glissando
		if s.voice.glissando != nil {
			t := 100 * s.voice.envelopePosition / int(sampleRate)
			s.voice.note = s.voice.glissando[t]
		}
		// Apply envelope
		attackPosition := int((float64(s.voice.envelope.attackTime) / 100) * float64(sampleRate))
		decayPosition := attackPosition + int((float64(s.voice.envelope.decayTime)/100)*float64(sampleRate))
		sustainPosition := decayPosition + int((float64(s.voice.envelope.sustainTime)/100)*float64(sampleRate))
		releasePosition := sustainPosition + int((float64(s.voice.envelope.releaseTime)/100)*float64(sampleRate))
		if s.voice.envelopePosition >= 0 && s.voice.envelopePosition <= attackPosition {
			attackRamp := float64(s.voice.envelope.attackLevel) / float64(attackPosition)
			currentLevel := int(attackRamp * float64(s.voice.envelopePosition))
			b = int16(w * float64(currentLevel) * float64(amplitudeScale))
			s.voice.released = false
			s.voice.envelopePosition++
		}
		if s.voice.envelopePosition > attackPosition && s.voice.envelopePosition <= decayPosition {
			decayRamp := float64(s.voice.envelope.decayLevel-s.voice.envelope.attackLevel) / (float64(decayPosition) - float64(attackPosition))
			currentLevel := int((decayRamp * float64(s.voice.envelopePosition-attackPosition)) + float64(s.voice.envelope.attackLevel))
			b = int16(w * float64(currentLevel) * float64(amplitudeScale))
			s.voice.released = false
			s.voice.envelopePosition++
		}
		if s.voice.envelopePosition > decayPosition && s.voice.envelopePosition <= sustainPosition {
			b = int16(w * float64(s.voice.envelope.sustainLevel) * float64(amplitudeScale))
			s.voice.released = false
			s.voice.envelopePosition++
		}
		if s.voice.envelopePosition > sustainPosition && s.voice.envelopePosition <= releasePosition {
			releaseRamp := float64(0-s.voice.envelope.sustainLevel) / (float64(releasePosition) - float64(sustainPosition))
			currentLevel := int((releaseRamp * float64(s.voice.envelopePosition-sustainPosition)) + float64(s.voice.envelope.sustainLevel))
			b = int16(w * float64(currentLevel) * float64(amplitudeScale))
			s.voice.released = true
			s.voice.envelopePosition++
		}
		if s.voice.envelopePosition > sustainPosition {
			s.voice.released = true
		}
		if s.voice.envelopePosition == releasePosition {
			b = 0
			s.voice.released = true
			s.voice.envelopePosition = -1
		}
		if s.voice.envelopePosition < 0 || s.voice.envelopePosition > releasePosition {
			s.voice.released = true
			b = 0
		}

		buf[4*i] = byte(b)
		buf[4*i+1] = byte(b >> 8)
		buf[4*i+2] = byte(b)
		buf[4*i+3] = byte(b >> 8)
		p++
	}

	s.position += int64(len(buf))
	if int64(math.Round(length*4)) >= 1 {
		s.position %= int64(length * 4)
	} else {
		// This results in massive aliasing above C5 but idk maybe the Nimbus wasn't much better?
		s.position %= 4
	}

	if origBuf != nil {
		n := copy(origBuf, buf)
		s.remaining = buf[n:]
		return n, nil
	}
	return len(buf), nil
}

// Close is io.Closer's Close.
func (s *stream) Close() error {
	return nil
}

// Queue returns the number of free slots in a voice queue
func (n *Nimbus) Queue(v int) int {
	n.voices[v].muQueue.Lock()
	retval := 6 - len(n.voices[v].queue)
	n.voices[v].muQueue.Unlock()
	return retval
}

// AskSound returns true if the sound engine is on
func (n *Nimbus) AskSound() bool {
	return n.sound
}

// AskVoice returns the selected voice
func (n *Nimbus) AskVoice() int {
	return n.selectedVoice
}

// SetSound turns the sound engine on and off
func (n *Nimbus) SetSound(v bool) {
	if n.sound {
		return // already on
	}
	if v {
		n.envelopes = make([]envelope, 0)
		for i := 0; i < 13; i++ {
			// Envelopes 10,11,12 are not accesible and are used when duration is given
			n.envelopes = append(n.envelopes, envelope{sustainTime: (i + 1) * 5, sustainLevel: 15})
		}
		n.voices = make([]voice, 3)
		for i := 0; i < 3; i++ {
			n.voices[i].queue = make([]noteItem, 0)
			n.voices[i].note = 1
			n.voices[i].envelope.attackTime = 0
			n.voices[i].envelope.attackLevel = 15
			n.voices[i].envelope.decayTime = 0
			n.voices[i].envelope.decayLevel = 15
			n.voices[i].envelope.sustainTime = 20
			n.voices[i].envelope.sustainLevel = 15
			n.voices[i].envelope.releaseTime = 0
			n.voices[i].envelopePosition = -1
			n.voices[i].selectedEnvelope = 1
		}
		n.sound = true
		n.selectedVoice = 0
		n.selectedEnvelope = 0
		go n.PlayQueues()
	} else {
		n.sound = false
	}
}

func (n *Nimbus) SetVoice(v int) {
	n.selectedVoice = v
}

func (n *Nimbus) SetEnvelope(e int) {
	n.selectedEnvelope = e
}

func (n *Nimbus) Note(pitch1, pitch2, duration, volume, envelope int) bool {
	var note noteItem
	note.note = pitch1
	note.note2 = pitch2
	note.volume = volume
	// Handle envelope vs volume
	if volume > 0 {
		// use passed envelope or preselected envelope if no option passed
		if envelope > 0 {
			note.envelope = n.envelopes[envelope]
		} else {
			note.envelope = n.envelopes[n.voices[n.selectedVoice].selectedEnvelope]
		}
	}
	// Handle glissando
	if pitch2 > 0 {
		note.glissando = make(map[int]int)
		for i := 0; i < duration; i++ {
			note.glissando[i] = pitch1 + int((float64(pitch2)-float64(pitch1))*(float64(i)/float64(duration)))
		}
	}
	// If we've got a duration then used envelope 10 with sustain time set to duration
	if duration > 0 {
		n.envelopes[10+n.selectedVoice].sustainTime = duration
		n.envelopes[10+n.selectedVoice].sustainLevel = 15
		note.envelope = n.envelopes[10+n.selectedVoice]
	} else {
		note.envelope = n.envelopes[n.selectedEnvelope]
	}
	// Handle volume
	if volume > 0 {
		note.envelope.sustainLevel = volume
	}
	return n.voices[n.selectedVoice].Add(note)
}

// Bell makes the Nimbus bell sound
func (n *Nimbus) Bell() {
	// Store current voice settings so we can revert back
	oldnote := n.voices[0].note
	oldattackTime := n.voices[0].envelope.attackTime
	oldattackLevel := n.voices[0].envelope.attackLevel
	olddecayTime := n.voices[0].envelope.decayTime
	olddecayLevel := n.voices[0].envelope.decayLevel
	oldsustainTime := n.voices[0].envelope.sustainTime
	oldsustainLevel := n.voices[0].envelope.sustainLevel
	oldreleaseTime := n.voices[0].envelope.releaseTime

	// Ring da bell
	n.voices[0].note = 148
	n.voices[0].envelope.attackTime = 0
	n.voices[0].envelope.attackLevel = 15
	n.voices[0].envelope.decayTime = 0
	n.voices[0].envelope.decayLevel = 15
	n.voices[0].envelope.sustainTime = 20
	n.voices[0].envelope.sustainLevel = 15
	n.voices[0].envelope.releaseTime = 0
	n.voices[0].Trigger()
	for !n.voices[0].released {
		time.Sleep(10 * time.Millisecond)
	}

	// Revert settings
	n.voices[0].note = oldnote
	n.voices[0].envelope.attackTime = oldattackTime
	n.voices[0].envelope.attackLevel = oldattackLevel
	n.voices[0].envelope.decayTime = olddecayTime
	n.voices[0].envelope.decayLevel = olddecayLevel
	n.voices[0].envelope.sustainTime = oldsustainTime
	n.voices[0].envelope.sustainLevel = oldsustainLevel
	n.voices[0].envelope.releaseTime = oldreleaseTime
}

// PlayQueue plays all the notes in all the queues
func (n *Nimbus) PlayQueues() {
	for n.sound {
		for v := 0; v < len(n.voices); v++ {
			if len(n.voices[v].queue) > 0 {
				go n.voices[v].Play(n)
			}
		}
		time.Sleep(10 * time.Millisecond)
	}
}
