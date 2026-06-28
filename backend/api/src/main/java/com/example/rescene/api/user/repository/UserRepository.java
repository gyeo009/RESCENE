package com.example.rescene.api.user.repository;

import com.example.rescene.api.user.domain.UserEntity;
import org.springframework.data.jpa.repository.JpaRepository;

public interface UserRepository extends JpaRepository<UserEntity, Long> {
}
