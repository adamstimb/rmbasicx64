package nimgobus

import (
	"log"
	"math"
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten/v2/audio"
)

// Nimbus sound synthesis!

const (
	sampleRate   = 44100
	baseFreq     = 55 //220
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
	note     int
	envelope envelope
}

type voice struct {
	audioContext     *audio.Context
	player           *audio.Player
	note             int
	envelope         envelope
	envelopePosition int
	released         bool
	muQueue          sync.Mutex
	queue            []noteItem
}

// Add a note to the queue if the queue is not full
func (v *voice) Add(item noteItem) bool {
	if len(v.queue) >= 6 {
		return false
	} else {
		v.muQueue.Lock()
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
		v.note = item.note
		v.envelope.attackTime = item.envelope.attackTime
		v.envelope.attackLevel = item.envelope.attackLevel
		v.envelope.decayTime = item.envelope.decayTime
		v.envelope.decayLevel = item.envelope.decayLevel
		v.envelope.sustainTime = item.envelope.sustainTime
		v.envelope.sustainLevel = item.envelope.sustainLevel
		v.envelope.releaseTime = item.envelope.releaseTime
		v.Trigger()
		for !v.released {
			time.Sleep(10 * time.Millisecond)
		}
		v.queue = v.queue[1:]
		v.muQueue.Unlock()
	} else {
		v.muQueue.Unlock()
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
	octave := int(math.Floor((float64(note)-1.0)/12.0)) + 1
	note = note + 9 // because lowest Nimbus note is C, not A
	freq := float64(baseFreq) * math.Pow(2, (float64(note)/12))
	return freq * float64(octave)
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

	length := int64(sampleRate / frequencyOfNote(s.voice.note))
	amplitudeScale := maxAmplitude / 16
	p := s.position / 4
	for i := 0; i < len(buf)/4; i++ {
		if !s.printOnce {
			s.printOnce = true
		}
		var b int16
		w := math.Sin(2 * math.Pi * float64(p) / float64(length))
		// Convert sine wave to square
		if w > 0 {
			w = 1
		} else {
			w = -1
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
	s.position %= length * 4

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

// SetSound turns the sound engine on and off
func (n *Nimbus) SetSound(v bool) {
	if v {
		n.envelopes = make([]envelope, 0)
		for i := 0; i < 10; i++ {
			n.envelopes = append(n.envelopes, envelope{sustainTime: (i + 1) * 5, sustainLevel: 15})
		}
		n.voice1.queue = make([]noteItem, 0)
		n.voice1.note = 1
		n.voice1.envelope.attackTime = 0
		n.voice1.envelope.attackLevel = 15
		n.voice1.envelope.decayTime = 0
		n.voice1.envelope.decayLevel = 15
		n.voice1.envelope.sustainTime = 20
		n.voice1.envelope.sustainLevel = 15
		n.voice1.envelope.releaseTime = 0
		n.voice1.envelopePosition = -1
		n.voice2.queue = make([]noteItem, 0)
		n.voice2.note = 1
		n.voice2.envelope.attackTime = 0
		n.voice2.envelope.attackLevel = 15
		n.voice2.envelope.decayTime = 0
		n.voice2.envelope.decayLevel = 15
		n.voice2.envelope.sustainTime = 20
		n.voice2.envelope.sustainLevel = 15
		n.voice2.envelope.releaseTime = 0
		n.voice2.envelopePosition = -1
		n.sound = true
		go n.PlayQueues()
	} else {
		n.sound = false
	}
}

func (n *Nimbus) TestSound() {
	log.Printf("Adding notes")
	n.voice1.Add(noteItem{note: 18, envelope: n.envelopes[9]})
	n.voice1.Add(noteItem{note: 12, envelope: n.envelopes[9]})
	n.voice1.Add(noteItem{note: 6, envelope: n.envelopes[9]})
	n.voice1.Add(noteItem{note: 1, envelope: n.envelopes[9]})
	n.voice2.Add(noteItem{note: 18, envelope: n.envelopes[4]})
	n.voice2.Add(noteItem{note: 12, envelope: n.envelopes[4]})
	n.voice2.Add(noteItem{note: 6, envelope: n.envelopes[4]})
	n.voice2.Add(noteItem{note: 1, envelope: n.envelopes[4]})
	log.Printf("Adding notes done")
}

// Bell makes the Nimbus bell sound
func (n *Nimbus) Bell() {
	// Store current voice settings so we can revert back
	oldnote := n.voice1.note
	oldattackTime := n.voice1.envelope.attackTime
	oldattackLevel := n.voice1.envelope.attackLevel
	olddecayTime := n.voice1.envelope.decayTime
	olddecayLevel := n.voice1.envelope.decayLevel
	oldsustainTime := n.voice1.envelope.sustainTime
	oldsustainLevel := n.voice1.envelope.sustainLevel
	oldreleaseTime := n.voice1.envelope.releaseTime

	// Ring da bell
	n.voice1.note = 9
	n.voice1.envelope.attackTime = 0
	n.voice1.envelope.attackLevel = 15
	n.voice1.envelope.decayTime = 0
	n.voice1.envelope.decayLevel = 15
	n.voice1.envelope.sustainTime = 20
	n.voice1.envelope.sustainLevel = 15
	n.voice1.envelope.releaseTime = 0
	n.voice1.Trigger()
	for !n.voice1.released {
		time.Sleep(10 * time.Millisecond)
	}

	// Revert settings
	n.voice1.note = oldnote
	n.voice1.envelope.attackTime = oldattackTime
	n.voice1.envelope.attackLevel = oldattackLevel
	n.voice1.envelope.decayTime = olddecayTime
	n.voice1.envelope.decayLevel = olddecayLevel
	n.voice1.envelope.sustainTime = oldsustainTime
	n.voice1.envelope.sustainLevel = oldsustainLevel
	n.voice1.envelope.releaseTime = oldreleaseTime
}

// PlayQueue plays all the notes in all the queues
func (n *Nimbus) PlayQueues() {
	for n.sound {
		if len(n.voice1.queue) > 0 {
			go n.voice1.Play(n)
		}
		if len(n.voice2.queue) > 0 {
			go n.voice2.Play(n)
		}
		if len(n.voice2.queue) > 0 {
			go n.voice2.Play(n)
		}
		time.Sleep(10 * time.Millisecond)
	}
}
