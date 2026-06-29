package com.example.rescene.api.user.domain;

import com.example.rescene.common.jpa.entity.AuditableEntity;
import jakarta.persistence.Column;
import jakarta.persistence.Entity;
import jakarta.persistence.GeneratedValue;
import jakarta.persistence.GenerationType;
import jakarta.persistence.Id;
import jakarta.persistence.Table;

@Entity
@Table(
        name = "users"
)
public class UserEntity extends AuditableEntity {
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;

	@Column(name = "display_name", nullable = false, length = 50)
	private String displayName;

    public Long getId() {
        return id;
    }

    public String getDisplayName(){
        return displayName;
    }

}
