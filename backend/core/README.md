Entity에 @Setter, @Data, @AllArgsConstructor 금지
JPA용 기본 생성자는 protected
생성은 create(...) 정적 팩토리
상태 변경은 의미 있는 메서드
도메인 규칙은 도메인 테스트로 먼저 고정
Service는 repository/transaction/use case 조합 담당
Controller는 HTTP 계약만 담당
common은 진짜 공통 기반만
