/PROG  GET_ZONE
/ATTR
DEFAULT_GROUP	= 1,*,*,*,*;
/MN
 : CALL AT_PICK ;
 : IF R{zoneID}=${ZONE_PICK},JMP LBL[999] ;
 :  ;
 : CALL AT_PLACE ;
 : IF R{zoneID}=${ZONE_PLACE},JMP LBL[999] ;
 :  ;
 : CALL AT_HOME ;
 : IF R{zoneID}=${ZONE_HOME},JMP LBL[999] ;
 :  ;
 : R{zoneID}=${ZONE_UNKNOWN} ;
 : END ;
 :  ;
 : LBL[999] ;
 :   ! zone found ;
/POS
/END
