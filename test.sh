#!/bin/bash

mkdir /tmp/dirwatcher

processes=()
dirs=()

for i in {1..10}
do
	dirs+=("/tmp/dirwatcher/dir_$i")
	mkdir /tmp/dirwatcher/dir_$i
	dirwatcher /tmp/dirwatcher/dir_$i &
	processes+=($!)
done

echo "Processes started, press any key to continue"
read -n 1

for dir in "${dirs[@]}"
do
	for ext in {1..100}
	do
		echo "test" > ${dir}/test.${ext}
	done
done

echo "Files created, press any key to continue"
read -n 1

for dir in "${dirs[@]}"
do
	for ext in {1..10}
	do
		if (($ext % 2)); then
			echo "test" >> ${dir}/test.${ext}
		fi
		# if (($ext % 3)); then
		# 	mv ${dir}/test.${ext} ${dir}/test.${ext}.OLD
		# fi
		if (($ext % 7)); then
			rm -f ${dir}/test.${ext}
			# rm -f ${dir}/test.${ext}.OLD
		fi
	done
done

echo "Files manipulated, press any key to continue"
read -n 1

for i in "${processes[@]}"
do
	kill -INT $i
done

rm -rf /tmp/dirwatcher
echo "All cleaned up!"
