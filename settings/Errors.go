package settings;

import (
    "github.com/barbell-math/block/util/customerr"
)

var DataVersionMalformed,IsDataVersionMalformed=customerr.ErrorFactory(
    "Data version is malformed.",
);

var SettingsFileNotFound,IsSettingsFileNotFound=customerr.ErrorFactory(
    "A file specified in the settings config file could not be found.",
);
