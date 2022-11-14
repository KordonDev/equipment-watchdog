package members

type MemberService struct {
	memberDB *memberDB
}

func NewMemberService(memberDB *memberDB) *MemberService {
	return &MemberService{
		memberDB: memberDB}

}
