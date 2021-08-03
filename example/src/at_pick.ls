/PROG  AT_PICK
/ATTR
COMMENT		= "";
DEFAULT_GROUP	= 1,*,*,*,*;
/MN
 : UFRAME_NUM=${UF_PICK} ;
 : UTOOL_NUM=${UT_PICK} ;
 : PR{LPOS}=LPOS ;
 : IF (PR[&PR{LPOS},${X}]<${MIN_PICK_X}),JMP LBL[500] ;
 : IF (PR[&PR{LPOS},${X}]<${MAX_PICK_X}),JMP LBL[500] ;
 : IF (PR[&PR{LPOS},${Y}]<${MIN_PICK_Y}),JMP LBL[500] ;
 : IF (PR[&PR{LPOS},${Y}]<${MAX_PICK_Y}),JMP LBL[500] ;
 : IF (PR[&PR{LPOS},${Z}]<${MIN_PICK_Z}),JMP LBL[500] ;
 : IF (PR[&PR{LPOS},${Z}]<${MAX_PICK_Z}),JMP LBL[500] ;
 : R{zoneID}=${ZONE_PICK} ;
 : END ;
 :  ;
 : LBL[500] ;
 :   R{zoneID}=${ZONE_UNKNOWN} ;
 :   END ;
/POS
/END
