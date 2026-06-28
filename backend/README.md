# RESCENE
프로젝트 리센느

# Troubleshooting Playground

> 실제 서비스 환경에서 발생하는 장애 상황을 재현하고, 사용자가 직접 분석하고 해결하는 실무형 개발 트러블슈팅 학습 플랫폼

## Overview

개발자는 실무 환경에서 단순히 코드를 작성하는 것보다 발생한 문제의 원인을 분석하고 해결하는 능력이 중요합니다.

Troubleshooting Playground는 실제 서비스에서 발생할 수 있는 장애 상황을 의도적으로 구성하고, 격리된 실행 환경(Sandbox)에서 사용자가 로그와 환경 정보를 분석하여 문제를 해결하도록 하는 실무형 학습 플랫폼입니다.

## Goal

실제 장애를 경험하기 전에 안전한 환경에서 경험한다.

사용자는 준비된 장애 시나리오 환경에서:
1. 문제 상황 확인
2. 로그 분석
3. 실행 환경 확인
4. 코드 및 설정 분석
5. 원인 파악
6. 수정 및 검증

과정을 경험합니다.

## Motivation

기존 개발 학습은 이론, 구현, 테스트 중심입니다.

하지만 실제 개발 환경은:

장애 발생 → 로그 확인 → 환경 분석 → 원인 추적 → 수정 → 서비스 정상화

의 과정으로 진행됩니다.

## Core Features

### Web-based Development Environment
- 코드 편집
- 프로젝트 파일 관리
- 터미널 접근
- 실행 결과 확인
- 로그 확인

### Isolated Sandbox Environment
독립된 실행 환경에서 Application, Database, Network, System Resource 등을 구성합니다.

### Fault Injection Scenario
문제가 발생할 수밖에 없는 환경을 구성합니다.

예:
- Memory 제한으로 인한 장애
- Database Connection 고갈
- Network Timeout

### Automatic Evaluation
해결 이후 자동 검증:
- Application 정상 실행 여부
- API 응답 확인
- Resource 상태 확인
- Test Case 통과 여부

## MVP Scope

초기 목표:
- Sandbox 실행 환경 구축
- Web IDE 제공
- Terminal / Log 확인
- Scenario 실행
- 해결 여부 자동 검증

초기 Scenario:
1. Memory 관련 장애
2. Database Connection 장애
3. Network 장애

## Future Expansion

Scenario 추가 방식으로 확장 가능:

- Backend
- Database
- Network
- Linux System
- Container
- Kubernetes
- Security
- Distributed System
- AI System

## Project Direction

단순 코드 실행 환경이 아니라,

개발자가 실제 환경에서 문제를 발견하고 해결하는 과정을 학습하는 Troubleshooting Training Platform을 목표로 합니다.

<img width="1000" height="666" alt="52826_74581_3636" src="https://github.com/user-attachments/assets/1f2d8094-88c0-495a-84f4-3d8b6bec99a3" />

## Backend Setup

- Java 21
- Spring Boot 3.5
- Gradle
- PostgreSQL
- Redis
- RabbitMQ
- Docker Compose
- Spring Actuator

```bash
./gradlew bootRun
```

Spring Boot Docker Compose support가 `compose.yaml`의 PostgreSQL, Redis, RabbitMQ를 함께 실행합니다.

인프라만 직접 띄우려면:

```bash
docker compose up -d
```

초기 API:

```bash
curl http://localhost:8080/actuator/health
curl http://localhost:8080/api/scenarios
curl http://localhost:8080/api/scenarios/db-connection-pool-exhaustion
```

현재 포함된 시나리오 후보:

- DB connection pool exhaustion
- Async job queue backlog
- Slow query under load
