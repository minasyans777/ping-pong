#!/bin/bash

read -p "Please write the full path of your program. " program_path

go build -o /home/sam/pingpong $program_path 
