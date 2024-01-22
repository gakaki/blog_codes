package com.gakaki.demo.service;

import org.junit.jupiter.api.Test;

import static com.gakaki.demo.service.SatigerBookService.regexFind;
import static org.junit.jupiter.api.Assertions.*;

class SatigerBookServiceTest {

    @Test
    void testRegexFind() {

        String regexStr = "your_regex_here";
        String inputStr = "your_regex_here121212";
        String result = regexFind(regexStr, inputStr);
        if (result != null) {
            System.out.println(result);
        } else {
            System.out.println("No match found.");
        }
    }
}