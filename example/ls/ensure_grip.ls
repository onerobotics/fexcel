/PROG  ENSURE_GRIP
/ATTR
DEFAULT_GROUP	= *,*,*,*,*;
/MN
 : LBL[1] ;
 :   WAIT (!RI[1:UNGRIPPED] AND RI[2:GRIPPED]) TIMEOUT,LBL[501] ;
 :   END ;
 :  ;
 : LBL[501] ;
 :   ! timeout ;
 :   ! TODO: throw error? ;
 :   JMP LBL[1] ;
/POS
/END
