package com.example.rescene.core.user.controller;

import com.example.rescene.core.user.dto.UserDisplayNameResponse;
import com.example.rescene.core.user.service.UserService;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RestController;

@RestController
public class UserController {
    private final UserService userService;

    public UserController(UserService userService) {
        this.userService = userService;
    }

    @GetMapping("/api/users/{id}/display-name")
    public ResponseEntity<UserDisplayNameResponse> getUserDisplayName(@PathVariable("id") Long id) {
        String displayName = userService.getUserDisplayName(id);
        return ResponseEntity.ok(new UserDisplayNameResponse(id, displayName));
    }
}
