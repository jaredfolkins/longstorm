#!/bin/bash
go get -u github.com/gobuffalo/packr/v2/packr2;
packr2 build; packr2 clean
