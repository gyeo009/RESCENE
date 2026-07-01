package com.example.rescene.core.user.repository;

import com.example.rescene.core.user.domain.UserEntity;
import org.springframework.data.jpa.repository.JpaRepository;

public interface UserRepository extends JpaRepository<UserEntity, Long> {
}
