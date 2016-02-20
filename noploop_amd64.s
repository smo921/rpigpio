TEXT ·shortWait(SB), $0
      MOVL    cnt+8(FP), CX
      MOVL    $0, AX
      CMPL    AX, CX
      JGE     L2
L1:   INCL    AX
      NOP
      CMPL    AX, CX
      JCS L1
L2:   RET
