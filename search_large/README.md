Одна из часто встречающихся задач в программах - поиск элемента в множестве. В случае когда у нас числа и они хранятся в массиве, стандартное решение - двоичный поиск. Во многих языках программирования, в том числе и C++, двоичный поиск реализован в библиотеке. В случае C++ - `std::find()`.

Это уже третья задача, в которой вам предлагается ускорить стандартную реалзацию. Подумайте, что может замедлять работу программы, когда данных много? Какие метрики могут подтвердить или опровергнуть вашу гипотезу? Попробуйте проанализировать их с помощью `perf`. Подумайте, что можно добавить к вашему решению задачи `search_medium`, чтобы ускорить работу в случае, когда объём данных больше мегабайта? Можете считать, что вам доступны расширения x86 вплоть до AVX2.

Реализуйте быстрый поиск по отсортированному массиву большого размера: от `2**16` до `2**20` элементов. Напишите своё решение в файле `search.cpp`.
