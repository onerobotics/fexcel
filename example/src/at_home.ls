/PROG  AT_HOME
/ATTR
COMMENT		= "";
DEFAULT_GROUP	= 1,*,*,*,*;
/MN
 : UFRAME_NUM=${WORLD} ;
 : UTOOL_NUM=${FACEPLATE} ;
 : PR{LPOS}=LPOS ;
 : IF (PR[&PR{LPOS},${X}]<${MIN_HOME_X}),JMP LBL[500] ;
 : IF (PR[&PR{LPOS},${X}]<${MAX_HOME_X}),JMP LBL[500] ;
 : IF (PR[&PR{LPOS},${Y}]<${MIN_HOME_Y}),JMP LBL[500] ;
 : IF (PR[&PR{LPOS},${Y}]<${MAX_HOME_Y}),JMP LBL[500] ;
 : IF (PR[&PR{LPOS},${Z}]<${MIN_HOME_Z}),JMP LBL[500] ;
 : IF (PR[&PR{LPOS},${Z}]<${MAX_HOME_Z}),JMP LBL[500] ;
 : R{zoneID}=${ZONE_HOME} ;
 : END ;
 :  ;
 : LBL[500] ;
 :   R{zoneID}=${ZONE_UNKNOWN} ;
 :   END ;
/POS
/END
