# -*- coding: utf-8 -*-

# Copyright (C) 2018 Nippon Telegraph and Telephone Corporation.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
# implied.
# See the License for the specific language governing permissions and
# limitations under the License.

import collections

"""
XML Element module
"""

from lxml import etree as ET

def element_name(name, space):
    """
    Create element name.
    """
    if space is None:
        return name

    return "{{{0}}}{1}".format(space, name)

def strbool(val):
    """
    convert bool to string.
    """
    return "true" if val else "false"

def _xmlstring(val):
    if isinstance(val, bool):
        return strbool(val)

    if isinstance(val, str):
        return val

    if val is not None:
        return str(val)

    return None

# pylint: disable=too-few-public-methods
class Element(object):
    """
    XML Element class.
    """
    _FIELDS = ()
    def __init__(self, name, text=None, space=None, nsmap=None):
        nsmap = nsmap if nsmap is not None else dict()
        if space is not None:
            nsmap[None] = space

        self._name = element_name(name, space)
        self._text = _xmlstring(text)
        self._nsmap = nsmap
        self._attrs = collections.OrderedDict()

    def xml_element(self):
        """
        Create Element.
        """
        elm = ET.Element(self._name, nsmap=self._nsmap, **self._attrs)
        if self._text is not None:
            elm.text = self._text

        return elm


class ListElement(Element):
    """
    XML list element class.
    """
    def __init__(self, name, *args, **kwargs):
        super(ListElement, self).__init__(name, *args, **kwargs)
        self._children = list()

    def append(self, *args):
        """
        append children.
        """
        for arg in args:
            self._children.append(arg)

    def xml_element(self):
        """
        create XML element.
        """
        elm = super(ListElement, self).xml_element()
        for child in self._children:
            elm.append(child.xml_element())
        return elm

    def merge_element(self, src):
        for child in src._children:
            self.append(child)


    def __len__(self):
        return len(self._children)


class DictElement(Element):
    """
    XML dict element class
    """
    def __init__(self, name, *args, **kwargs):
        super(DictElement, self).__init__(name, *args, **kwargs)
        self._children = collections.OrderedDict()

    def append(self, key, val):
        """
        append children. 
        """
        self._children[key] = val

    def appends(self, **kwargs):
        """
        append children. 
        """
        for key, val in d.items():
            self.append(key. val)

    def get(self, key):
        """
        get child
        """
        return self._children.get(key)

    def xml_element(self):
        """
        create XML element.
        """
        elm = super(DictElement, self).xml_element()
        for _, child in self._children.items():
            elm.append(child.xml_element())
        return elm

    def	merge_element(self, src):
        for key, srcval in src._children.items():
            if key not in self._children:
                self.append(key, srcval)
            else:
                self._children[key].merge_element(srcval)


def _to_elment_name(name):
    if name.startswith("_"):
        name = name[1:]
    return name.replace("_", "-").replace("--", "_")


class BaseElement(Element):
    """
    XML base element
    """

    def xml_element(self):
        elm = super(BaseElement, self).xml_element()
        for field_name in self._FIELDS:
            if not hasattr(self, field_name):
                continue
            field = getattr(self, field_name)
            if field is None:
                continue
            elif hasattr(field, "xml_element"):
                elm.append(field.xml_element())
            else:
                elm.append(xml_element(_to_elment_name(field_name), field))

        return elm

    def merge_element(self, src):
        merge_elements(self, src, *self._FIELDS)


def xml_element(name, text, space=None, nsmap=None):
    """
    Create lxml.Element.
    """
    return Element(name, text=text, space=space, nsmap=nsmap).xml_element()


def merge_element(dst, src, name):
    if src is None:
        return

    dst_elm = getattr(dst, name, None)
    src_elm = getattr(src, name, None)

    if dst_elm is not None and hasattr(dst_elm, "merge_element"):
        dst_elm.merge_element(src_elm)

    elif src_elm is not None:
        setattr(dst, name, src_elm)


def merge_elements(dst, src, *names):
    for name in names:
        merge_element(dst, src, name)
