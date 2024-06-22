#include <thread>
#include <mutex>
#include <string>
#include <unordered_map>
#include <iostream>
#include <vector>

std::unordered_map<std::string, int> map;
std::mutex map_mutex;

void inc(const std::string& key) {
    auto guard = std::lock_guard(map_mutex);
    map[key]++;
}

void func() {
    for (int i = 0; i < 100000000; ++i) {
        inc(std::to_string(i % 100000));
    }
}

void do_collect(int *result) {
    auto guard = std::lock_guard(map_mutex);
    for (auto [_, v] : map) {
        *result += v;
    }
}
void collect(int *result) {
    while (true) {
        do_collect(result);
    }
}

int main() {
    std::vector<std::thread> threads;
    int result;
    threads.emplace_back(collect, &result);
    for (int i = 0; i < 2; ++i) {
        threads.emplace_back(func);
    }
    for (auto& thread : threads) {
        thread.join();
    }
}
