# -*- Autoconf -*-
# Process this file with autoconf to produce a configure script.

AC_PREREQ([2.69])
AC_INIT([GoNetC], [0.1], [admin@example.net])
AM_INIT_AUTOMAKE([foreign])

# Checks for programs.
AC_PROG_INSTALL

# Checks for libraries.

# Checks for header files.

# Checks for typedefs, structures, and compiler characteristics.

# Checks for library functions.

FF_HOME=`pwd`
AC_SUBST([FF_HOME])

AC_ARG_WITH([go-prefix],
            [AS_HELP_STRING([--with-go-prefix],
	    [GOPATH @<:@default=~/go@:>@])],
	    [],
	    [with_go_prefix=~/go])

GOPREFIX="$with_go_prefix"
AC_SUBST([GOPREFIX])

#
# Makefiles
#
AC_CONFIG_FILES([Makefile
                 src/Makefile
		 src/netconf/Makefile
		 src/netconf/lib/Makefile
		 src/netconf/lib/gobgp/Makefile
		 src/netconf/lib/gobgp/openconfig/Makefile
		 src/netconf/lib/lxd/Makefile
		 src/netconf/lib/net/Makefile
		 src/netconf/lib/netplan/Makefile
		 src/netconf/lib/openconfig/Makefile
		 src/netconf/lib/property/Makefile
		 src/netconf/lib/signal/Makefile
		 src/netconf/lib/sysctl/Makefile
		 src/netconf/lib/sysrepo/Makefile
		 src/netconf/lib/vty/Makefile
		 src/netconf/lib/xml/Makefile
		 src/netconf/app/Makefile
		 src/netconf/app/ncm/Makefile
		 src/netconf/app/ncm/dbm/Makefile
		 src/netconf/app/ncm/modules/Makefile
		 src/netconf/app/ncm/scripts/Makefile
		 src/netconf/app/cfg/Makefile
		 src/netconf/app/cfg/api/Makefile
	])
AC_OUTPUT
