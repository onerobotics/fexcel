/PROG  MV_RETREAT_HOME
/ATTR
DEFAULT_GROUP	= 1,*,*,*,*;
/MN
 : CALL GET_ZONE ;
 : IF R[2:zoneID]<>3,CALL UNSAFE ;
 :  ;
 : ! NOOP ;
/POS
/END
