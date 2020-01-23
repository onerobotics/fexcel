/PROG  TEST
/ATTR
COMMENT		= "";
DEFAULT_GROUP	= 1,*,*,*,*;
/APPL
AUTO_SINGULARITY_HEADER;
  ENABLE_SINGULARITY_AVOIDANCE   : TRUE;
/MN
 : ! this is a valid file ;
 : R{one}=R{two}+R{three} ;
 : J PR{home} 100% CNT100 ;
 : PR{lpos}=LPOS ;
 : PR{jpos}=JPOS ;
/POS
/END
