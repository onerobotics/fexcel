/PROG  ENSURE_UNGRIP
/ATTR
DEFAULT_GROUP	= *,*,*,*,*;
/MN
 : LBL[1] ;
 :   WAIT (RI{UNGRIPPED} AND !RI{GRIPPED}) TIMEOUT,LBL[501] ;
 :   END ;
 :  ;
 : LBL[501] ;
 :   ! timeout ;
 :   ! TODO: throw error? ;
 :   JMP LBL[1] ;
/POS
/END
