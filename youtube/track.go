/*
MIT License

Copyright (c) 2019 Oded Shapira

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

// package youtube abstracts searching and playing audio through youtube.
// currently only the audio/webm container with the opus codec is supported
package youtube

import (
	"errors"
	"github.com/jeffallen/seekinghttp"
	"lionPlayer/core"
	"lionPlayer/webm"
	"time"
)

type Track struct {
	VideoId  string
	Title    string
	Author   string
	Duration time.Duration
	IsStream bool
	source   *Source
	Format   *Format
}

// Return a playable of this track that can be played.
func (t Track) GetPlayable() (core.Playable, error) {
	return t.GetPlaySeekable()
}

// Return a playseekable of this track that can be played.
func (t Track) GetPlaySeekable() (core.PlaySeekable, error) {
	vurl, err := t.Format.GetValidUrl()
	if err != nil {
		return nil, err
	}

	res := seekinghttp.New(vurl)
	res.Client = &t.source.Client

	if size, err := res.Size(); err != nil {
		return nil, err
	} else if size == 0 {
		return nil, errors.New("got an empty request")
	}

	parser, err := webm.NewParser(res)

	if err != nil {
		return nil, err
	}

	file, err := parser.Parse()
	return file, nil
}
