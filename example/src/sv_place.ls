/PROG  SV_PLACE
/ATTR
DEFAULT_GROUP	= 1,*,*,*,*;
/MN
 : CALL MV_TO_PLACE ;
 :  ;
 : LBL[1] ;
 :   WAIT (DI{OK_TO_PLACE}) TIMEOUT,LBL[501] ;
 :   CALL ENSURE_GRIP ;
 :   CALL MV_PLACE ;
 :   CALL UNGRIP ;
 :   CALL MV_RETREAT_PLACE ;
 :   CALL ENSURE_UNGRIP ;
 :   END ;
 :  ;
 : LBL[501] ;
 :   ! TIMEOUT ;
 :   ! TODO: throw error ;
 :   JMP LBL[1] ;
/POS
/END
