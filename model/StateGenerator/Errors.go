package stateGenerator;

import (
    customerr "github.com/barbell-math/block/util/err"
)

var InvalidPredictionState,IsInvalidPredictionState=customerr.ErrorFactory(
    "Predictions cannot be made with the supplied parameters.",
);

var NoDataInSelectedWindow,IsNoDataInSelectedWindow=customerr.ErrorFactory(
    "No data was available within the selected window.",
);

var NoDataInSelectedTimeFrame,IsNoDataInSelectedTimeFrame=customerr.ErrorFactory(
    "No data was available within the selected time frame.",
);

var NotEnoughData,IsNotEnoughData=customerr.ErrorFactory(
    "Not enough data was available.",
);
