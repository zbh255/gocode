//go:build amd64

#include "textflag.h"

TEXT ·pause(SB),NOSPLIT,$0
    PAUSE
    RET
