SUBDIRS = netconf

clean-local:
	-rm -fr ./dist ./build
	-rm -fr ./beluganos_netconf.egg-info

dist-python:
	python ./setup.py sdist

install-python:
	python ./setup.py install
