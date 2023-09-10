package err;

var ValOutsideRange,IsValOutsideRange=ErrorFactory(
    "The specified value is outside the allowed range.",
);

var DimensionsDoNotAgree,IsDimensionsDoNotAgree=ErrorFactory(
    "Dimensions do not agree.",
);

var InvalidValue,IsInvalidValue=ErrorFactory(
    "The supplied value is not valid in the supplied context.",
);
