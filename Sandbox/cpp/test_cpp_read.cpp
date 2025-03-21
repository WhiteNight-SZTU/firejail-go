#include <iostream>

char conten[] = {
#include "fake.answer"
};

// char conte_1[] = {
// #include "true.answer"
// };

int main() {
    std::cout << conten << std::endl;
    std::cout << "Hello, World!" << std::endl;
    FILE *file = fopen("/home/whitenight/firejail-go/Sandbox/cpp/fake.answer", "r");
    if (file == NULL) {
        std::cout << "file open failed" << std::endl;
        return 1;
    }
    char buffer[1024];
    while (fgets(buffer, 1024, file) != NULL) {
        std::cout << buffer;
    }
    fclose(file);
    return 0;
}
