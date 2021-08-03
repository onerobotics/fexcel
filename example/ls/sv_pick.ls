/PROG  SV_PICK
/ATTR
DEFAULT_GROUP	= 1,*,*,*,*;
/MN
 : CALL MV_TO_PICK ;
 :  ;
 : LBL[1] ;
 :   WAIT (DI[1:OK_TO_PICK]) TIMEOUT,LBL[501] ;
 :   CALL ENSURE_UNGRIP ;
 :   CALL MV_PICK ;
 :   CALL GRIP ;
 :   CALL MV_RETREAT_PICK ;
 :   CALL ENSURE_GRIP ;
 :   END ;
 :  ;
 : LBL[501] ;
 :   ! TIMEOUT ;
 :   ! TODO: throw error ;
 :   JMP LBL[1] ;
/POS
/END
