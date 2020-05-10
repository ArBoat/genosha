package collyUnit

import (
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/debug"
)

func CollyInit()  {
	c := colly.NewCollector()
	c1 := colly.NewCollector(colly.Debugger(&debug.LogDebugger{}),
		colly.AllowURLRevisit(),
		)
}