/PROG  MV_RETREAT_PICK
/ATTR
DEFAULT_GROUP	= 1,*,*,*,*;
/MN
 : CALL GET_ZONE ;
 : IF R{zoneID}<>${ZONE_PICK},CALL UNSAFE ;
 :  ;
 : UFRAME_NUM=${UF_PICK} ;
 : UTOOL_NUM=${UT_PICK} ;
 :  ;
 : PR{LPOS}=LPOS ;
 : IF (PR[&PR{LPOS},${Z}]>${PICK_RETREAT_Z}),JMP LBL[1] ;
 :  ;
 : PR[&PR{LPOS},${Z}]=${PICK_RETREAT_Z} ;
 : L PR{LPOS} ${RETREAT_SPEED}mm/sec CNT0 ;
 :  ;
 : LBL[1] ;
 : J PR{PICK_PERCH} ${PERCH_SPEED}% CNT100 ;
/POS
/END
