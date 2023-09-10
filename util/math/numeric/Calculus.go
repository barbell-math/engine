package numeric

import (
	customerr "github.com/barbell-math/block/util/err"
	"github.com/barbell-math/block/util/math"
)

func Derivative[N math.SignedNumer](f func(x N) N, h N) func(x N) N {
	return func(x N) N {
        //return (f(x+h)-f(x-h))/(2*h);
        return (-f(x+2*h)+8*f(x+h)-8*f(x-h)+f(x-2*h))/(12*h);
	}
}

func Integral[N math.SignedNumer](
    f func(x N) N, 
) func(start N, end N, numPnts uint) (N,error) {
    return func(start N, end N, numPnts uint) (N,error) {
        if end<=start {
            return N(0),customerr.InvalidValue("End must be >start to run integration.");
        }
        if numPnts%2==0 || numPnts<3 {
            return N(0),customerr.InvalidValue("NumPnts must be an odd value >=3.");
        }
        rv:=N(0);
        //deltaX:=(end-start)/N(numPnts-1);
        for i:=uint(0); i<numPnts; i++ {
            x_n:=(end-start)*N(i)/N(numPnts-1)+start;
            if i==0 || i+1==numPnts {
                rv+=f(x_n);
            } else if i%2==0 {
                rv+=2*f(x_n);
            } else {
                rv+=4*f(x_n);
            }
        }
        return (end-start)/(N(numPnts-1)*N(3))*rv,nil;
    }
}

func ConstIntegralBound[N math.SignedNumer](c N) func (x N) N {
    return func(x N) N { return c; }
}
func DoubleIntegral[N math.SignedNumer](
    f func(x1 N, x2 N) N, 
) func(startX1 N, endX1 N, startX2 func(x N) N, endX2 func(x N) N, numPnts uint) (N,error) {
    return func(startX1, endX1 N, startX2, endX2 func(x N) N, numPnts uint) (N, error) {
        if endX1<=startX1 {
            return N(0),customerr.InvalidValue("EndX1 must be >startX1 to run integration.");
        }
        if numPnts%2==0 || numPnts<3 {
            return N(0),customerr.InvalidValue("NumPnts must be an odd value >=3.");
        }
        rv:=N(0);
        for i:=uint(0); i<numPnts; i++ {
            tmp:=N(0);
            x1_n:=(endX1-startX1)*N(i)/N(numPnts-1)+startX1;
            _endX2:=endX2(N(x1_n));
            _startX2:=startX2(N(x1_n));
            for j:=uint(0); j<numPnts; j++ {
                x2_n:=(_endX2-_startX2)*N(j)/N(numPnts-1)+_startX2;
                if j==0 || j+1==numPnts {
                    tmp+=f(x1_n,x2_n);
                } else if j%2==0 {
                    tmp+=2*f(x1_n,x2_n);
                } else {
                    tmp+=4*f(x1_n,x2_n);
                }
            }
            tmp=(_endX2-_startX2)/(N(numPnts-1)*N(3))*tmp;
            if i==0 || i+1==numPnts {
                rv+=tmp
            } else if i%2==0 {
                rv+=2*tmp;
            } else {
                rv+=4*tmp;
            }
        }
        return (endX1-startX1)/(N(numPnts-1)*N(3))*rv,nil;
    }
}
