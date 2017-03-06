package ishell

import (
	"fmt"
	"io"
	"sync"
	"time"
	"unicode/utf8"
)

// ProgressDisplay handles the display string for
// a progress bar.
type ProgressDisplay interface {
	// Determinate returns the strings to display
	// for percents 0 to 100.
	Determinate() [101]string
	// Indeterminate returns the strings to display
	// at interval.
	Indeterminate() []string
}

// ProgressBar is an ishell progress bar.
type ProgressBar interface {
	// Display sets the display of the progress bar.
	Display(ProgressDisplay)
	// Indeterminate sets the progress bar type
	// to indeterminate if true or determinate otherwise.
	Indeterminate(bool)
	// Interval sets the time between transitions for indeterminate
	// progress bar.
	Interval(time.Duration)
	// SetProgress sets the progress stage of the progress bar.
	// percent is from between 1 and 100.
	Progress(percent int)
	// Prefix sets the prefix for the output. The text to place before
	// the display.
	Prefix(string)
	// Suffix sets the suffix for the output. The text to place after
	// the display.
	Suffix(string)
	// Final sets the string to show after the progress bar is done.
	Final(string)
	// Start starts the progress bar.
	Start()
	// Stop stops the progress bar.
	Stop()
}

const progressInterval = time.Millisecond * 100

type progressBarImpl struct {
	display       ProgressDisplay
	indeterminate bool
	interval      time.Duration
	iterator      iterator
	percent       int
	prefix        string
	suffix        string
	final         string
	writer        io.Writer
	writtenLen    int
	running       bool
	wait          chan struct{}
	wMutex        sync.Mutex
	sync.Mutex
}

func newProgressBar(s *Shell) ProgressBar {
	display := simpleProgressDisplay{}
	return &progressBarImpl{
		interval:      progressInterval,
		writer:        s.writer,
		display:       display,
		iterator:      &stringIterator{set: display.Indeterminate()},
		indeterminate: true,
	}
}

func (p *progressBarImpl) Display(display ProgressDisplay) {
	p.display = display
}

func (p *progressBarImpl) Indeterminate(b bool) {
	p.indeterminate = b
}

func (p *progressBarImpl) Interval(t time.Duration) {
	p.interval = t
}

func (p *progressBarImpl) Progress(percent int) {
	if percent < 0 {
		percent = 0
	} else if percent > 100 {
		percent = 100
	}
	p.percent = percent
	p.indeterminate = false
	p.refresh()
}

func (p *progressBarImpl) Prefix(prefix string) {
	p.prefix = prefix
}

func (p *progressBarImpl) Suffix(suffix string) {
	p.suffix = suffix
}

func (p *progressBarImpl) Final(s string) {
	p.final = s
}

func (p *progressBarImpl) write(s string) error {
	p.erase(p.writtenLen)
	p.writtenLen = utf8.RuneCountInString(s)
	_, err := p.writer.Write([]byte(s))
	return err
}

func (p *progressBarImpl) erase(n int) {
	for i := 0; i < n; i++ {
		p.writer.Write([]byte{'\b'})
	}
}

func (p *progressBarImpl) done() {
	p.wMutex.Lock()
	defer p.wMutex.Unlock()

	p.erase(p.writtenLen)
	fmt.Fprintln(p.writer, p.final)
}

func (p *progressBarImpl) output() string {
	p.Lock()
	defer p.Unlock()

	var display string
	if p.indeterminate {
		display = p.iterator.next()
	} else {
		display = p.display.Determinate()[p.percent]
	}
	return fmt.Sprintf("%s%s%s ", p.prefix, display, p.suffix)
}

func (p *progressBarImpl) refresh() {
	p.wMutex.Lock()
	defer p.wMutex.Unlock()

	p.write(p.output())
}

func (p *progressBarImpl) Start() {
	p.Lock()
	p.running = true
	p.wait = make(chan struct{})
	p.Unlock()

	go func() {
		for {
			var running, indeterminate bool
			p.Lock()
			running = p.running
			indeterminate = p.indeterminate
			p.Unlock()

			if !running {
				break
			}
			time.Sleep(p.interval)
			if indeterminate {
				p.refresh()
			}
		}
		p.done()
		close(p.wait)
	}()
}

func (p *progressBarImpl) Stop() {
	p.Lock()
	p.running = false
	p.Unlock()

	<-p.wait
}

// ProgressDisplayCharSet is the character set for
// a progress bar.
type ProgressDisplayCharSet []string

// Determinate satisfies ProgressDisplay interface.
func (p ProgressDisplayCharSet) Determinate() [101]string {
	// TODO everything here works but not pleasing to the eyes
	// and probably not optimal.
	// This should be cleaner.
	var set [101]string
	for i := range set {
		set[i] = p[len(p)-1]
	}
	// assumption is than len(p) <= 101
	step := 101 / len(p)
	for i, j := 0, 0; i < len(set) && j < len(p); i, j = i+step, j+1 {
		for k := 0; k < step && i+k < len(set); k++ {
			set[i+k] = p[j]
		}
	}
	return set
}

// Indeterminate satisfies ProgressDisplay interface.
func (p ProgressDisplayCharSet) Indeterminate() []string {
	return p
}

// ProgressDisplayFunc is a convenience function to create a ProgressDisplay.
// percent is -1 for indeterminate and 0-100 for determinate.
type ProgressDisplayFunc func(percent int) string

// Determinate satisfies ProgressDisplay interface.
func (p ProgressDisplayFunc) Determinate() [101]string {
	var set [101]string
	for i := range set {
		set[i] = p(i)
	}
	return set
}

// Indeterminate satisfies ProgressDisplay interface.
func (p ProgressDisplayFunc) Indeterminate() []string {
	// loop through until we get back to the first string
	set := []string{p(-1)}
	for {
		next := p(-1)
		if next == set[0] {
			break
		}
		set = append(set, next)
	}
	return set
}

type iterator interface {
	next() string
}

type stringIterator struct {
	index int
	set   []string
}

func (s *stringIterator) next() string {
	current := s.set[s.index]
	s.index++
	if s.index >= len(s.set) {
		s.index = 0
	}
	return current
}

var (
	indeterminateCharSet = []string{
		"[====                ]", "[ ====               ]", "[  ====              ]",
		"[   ====             ]", "[    ====            ]", "[     ====           ]",
		"[      ====          ]", "[       ====         ]", "[        ====        ]",
		"[         ====       ]", "[          ====      ]", "[           ====     ]",
		"[            ====    ]", "[             ====   ]", "[              ====  ]",
		"[               ==== ]", "[                ====]",
		"[               ==== ]", "[              ====  ]", "[             ====   ]",
		"[            ====    ]", "[           ====     ]", "[          ====      ]",
		"[         ====       ]", "[        ====        ]", "[       ====         ]",
		"[      ====          ]", "[     ====           ]", "[    ====            ]",
		"[   ====             ]", "[  ====              ]", "[ ====               ]",
	}
	determinateCharSet = []string{
		"[                    ]", "[>                   ]", "[=>                  ]",
		"[==>                 ]", "[===>                ]", "[====>               ]",
		"[=====>              ]", "[======>             ]", "[=======>            ]",
		"[========>           ]", "[=========>          ]", "[==========>         ]",
		"[===========>        ]", "[============>       ]", "[=============>      ]",
		"[==============>     ]", "[===============>    ]", "[================>   ]",
		"[=================>  ]", "[==================> ]", "[===================>]",
	}
)

type simpleProgressDisplay struct{}

func (s simpleProgressDisplay) Determinate() [101]string {
	return ProgressDisplayCharSet(determinateCharSet).Determinate()
}
func (s simpleProgressDisplay) Indeterminate() []string {
	return indeterminateCharSet
}
