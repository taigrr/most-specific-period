# Most Specific Period

An MSP, or Most Specific Period is defined as follows: 
- Given a period that hasn’t started or has already finished, an error is returned.
- Given a period that is currently valid, the period which is valid is chosen.
- Given two valid periods, one which is a month long and one which is a week long and completely overlap, the week-long period is chosen.
- Given two valid periods, one which is a week long and one which is a month long, but only overlap by a day (in either direction) the week-long period is selected.
- Given two valid periods, each exactly 30 days long, but offset by a week, the second (newer) period is given precedence. 
- Given two valid periods with exact same start and end time, the period with the lowest-ranking lexicographic identifier is returned (i.e. period B has precedence over period A, as ‘B’ comes after ‘A’ lexicographically.) This is because the default behavior is to choose the newer period, and named periods are often in lexicographical order (increasing numbers, letters, etc.)

This library operates off Period interfaces, which contain the following:

```
type Period interface {
        GetStartTime() time.Time
        GetEndTime() time.Time
        GetIdentifier() string
}
```

An example program is available to observe the operation and usage of this library.
