package com.gakaki.demo.controller;

import com.gakaki.demo.model.SaltTigerBookItem;
import com.gakaki.demo.service.SatigerBookService;
import jakarta.annotation.Resource;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.GetMapping;

import java.util.List;

@Controller
public class BooksController
{
    @Resource
    SatigerBookService satigerBookService;

    @GetMapping("/books")
    public ResponseEntity<?> getBooks() {
        try {
            List<SaltTigerBookItem> items = satigerBookService.fetchBooks();
            return ResponseEntity.ok(items);
        } catch (Exception e) {
            return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR).body(e.getMessage());
        }
    }

    @GetMapping("/bookDetails")
    public ResponseEntity<?> getBookDetails() {
        try {
            List<SaltTigerBookItem> items = satigerBookService.fetchBookDetails();
            return ResponseEntity.ok(items);
        } catch (Exception e) {
            return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR).body(e.getMessage());
        }
    }
}
