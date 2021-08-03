/PROG  GET_ZONE
/ATTR
DEFAULT_GROUP	= 1,*,*,*,*;
/MN
 : CALL AT_PICK ;
 : IF R[2:zoneID]=1,JMP LBL[999] ;
 :  ;
 : CALL AT_PLACE ;
 : IF R[2:zoneID]=2,JMP LBL[999] ;
 :  ;
 : CALL AT_HOME ;
 : IF R[2:zoneID]=3,JMP LBL[999] ;
 :  ;
 : R[2:zoneID]=0 ;
 : END ;
 :  ;
 : LBL[999] ;
 :   ! zone found ;
/POS
/END
