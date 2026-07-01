package com.example.rescene.core.user.service;

import com.example.rescene.core.user.domain.UserEntity;
import com.example.rescene.core.user.repository.UserRepository;
import org.springframework.stereotype.Service;

@Service
public class UserService {
    private final UserRepository userRepository;

    public UserService(UserRepository userRepository) {
        this.userRepository = userRepository;
    }

    public String getUserDisplayName(Long id) {
        UserEntity user = userRepository.findById(id)
                .orElseThrow(() -> new IllegalArgumentException("User not found: " + id));

        return user.getDisplayName();
    }
}
