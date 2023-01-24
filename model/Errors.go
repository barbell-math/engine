package model;

import (
    customerr "github.com/barbell-math/block/util/err"
)

var ManyPredictions,IsManyPredictions=customerr.ErrorFactory(
    "Many predictions exist.",
);

var InvalidPredictionState,IsInvalidPredictionState=customerr.ErrorFactory(
    "Predictions cannot be made with the supplied parameters.",
);
