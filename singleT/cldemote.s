#include "go_asm.h"
#include "/home/ansible/xinran/goroot/src/runtime/go_tls.h"
#include "funcdata.h"
#include "textflag.h"
#include "/home/ansible/xinran/goroot/src/runtime/cgo/abi_amd64.h"

#define RARG1 AX

TEXT ·neg(SB), NOSPLIT, $0
    MOVQ     x+0(FP), AX
    NEGQ     AX
    MOVQ     AX, ret+8(FP)
    RET

TEXT ·cldemote(SB), NOSPLIT, $0-8
    MOVQ	addr+0(FP), AX
    CLDEMOTE (AX)
    RET
