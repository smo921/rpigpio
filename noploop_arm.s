TEXT ·shortWait(SB), $0
      MOVW    cnt(FP), R3
      MOVW    $0, R0
      MOVW    $0, R2
L1:   CMP     R3, R2
      BHS     $0, L2
      MOVW    $1, R1
      ADD     R1, R2, R2
      JMP     L1
L2:   JMP     (R14)
