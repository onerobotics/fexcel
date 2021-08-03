/PROG  AT_PLACE
/ATTR
COMMENT		= "";
DEFAULT_GROUP	= 1,*,*,*,*;
/MN
 : UFRAME_NUM=${UF_PLACE} ;
 : UTOOL_NUM=${UT_PLACE} ;
 : PR{LPOS}=LPOS ;
 : IF (PR[&PR{LPOS},${X}]<${MIN_PLACE_X}),JMP LBL[500] ;
 : IF (PR[&PR{LPOS},${X}]<${MAX_PLACE_X}),JMP LBL[500] ;
 : IF (PR[&PR{LPOS},${Y}]<${MIN_PLACE_Y}),JMP LBL[500] ;
 : IF (PR[&PR{LPOS},${Y}]<${MAX_PLACE_Y}),JMP LBL[500] ;
 : IF (PR[&PR{LPOS},${Z}]<${MIN_PLACE_Z}),JMP LBL[500] ;
 : IF (PR[&PR{LPOS},${Z}]<${MAX_PLACE_Z}),JMP LBL[500] ;
 : R{zoneID}=${ZONE_PLACE} ;
 : END ;
 :  ;
 : LBL[500] ;
 :   R{zoneID}=${ZONE_UNKNOWN} ;
 :   END ;
/POS
/END
