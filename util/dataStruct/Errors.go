package dataStruct;

import (
    customerr "github.com/barbell-math/engine/util/err"
)

var QueueFull,IsQueueFull=customerr.ErrorFactory(
    "The capacity of the queue has been reached.",
);

var QueueEmpty,IsQueueEmpty=customerr.ErrorFactory(
    "The queue is empty.",
);
