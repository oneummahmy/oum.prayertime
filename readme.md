# OU Prayer Time

The source of this prayer time is from Jakim Website as we scrape the page.

## Example Usage

Install package by running

```shell
go get -u github.com/oneummahmy/oum.prayertime
```

Example code

```go
package main

import (
	"fmt"
	"github.com/oneummahmy/oum.prayertime"
)

func main() {
	locations, err := oum_prayertime.GetLocations()
	if err != nil {
		panic(err)
	}

	fmt.Println(locations)

	prayer, err := oum_prayertime.GetPrayerTime(locations[0])
	if err != nil {
		panic(err)
	}

	fmt.Println(prayer.PrayerTime[0].Dhuhr) // the prayer time for Dhuhr is 12:50:00
}
```
