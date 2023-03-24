package log;

import (
    customerr "github.com/barbell-math/block/util/err"
)

var LogFileNotSpecified,IsLogFileNotSpecified=customerr.ErrorFactory(
    "The log file was not specified.",
);

var LogLineMalformed,IsLogLineMalformed=customerr.ErrorFactory(
    "Log line malformed.",
);
