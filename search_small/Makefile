CPP=g++
run: main
	./main
main: main.cpp search.cpp
	$(CPP) -g -march=native -O3 -o main main.cpp search.cpp
