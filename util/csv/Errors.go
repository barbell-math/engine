package csv;

import (
    "github.com/barbell-math/block/util/customerr"
)

var MalformedCSVFile,IsMalformedCSVFile=customerr.ErrorFactory(
    "The CSV file cannot be converted to the requested struct.",
);

var NonStructValue,IsNonStructValue=customerr.ErrorFactory(
    "A struct value was expected but was not recieved.",
);

