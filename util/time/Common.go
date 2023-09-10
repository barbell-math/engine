package time;

import (
    "time"
)

//After is the 'future' time, before is the 'past' time
func Between(after time.Time, before time.Time) (func(t time.Time) bool) {
    //if the 'past' time is after the 'future' time then switch them
    if before.After(after) {
        tmp:=before;
        before=after;
        after=tmp;
    }
    return func(t time.Time) bool {
        return (before.AddDate(0, 0, -1).Before(t) &&
            after.AddDate(0, 0, 1).After(t));
    }
}

func DaysBetween(after time.Time, before time.Time) int {
    if before.After(after) {
        tmp:=before;
        before=after;
        after=tmp;
    }
    return int(after.Sub(before).Hours()/24);
}
