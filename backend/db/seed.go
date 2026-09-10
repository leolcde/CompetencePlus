package db

import (
	"log"
	"strings"
	"time"

	"profilsactifs/models"

	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SeedQuestions(g *gorm.DB) error {
	var count int64
	if err := g.Model(&models.Question{}).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	yesOrNo := pq.StringArray{"Oui", "Non"}
	contents := []string{
		"I communicate easily with my colleagues.",
		"I am comfortable speaking in public.",
		"I can express myself clearly in writing.",
		"I ask questions when I do not understand.",
		"I express myself with confidence.",
		"I stay calm when facing disagreement.",
		"I respect agreed schedules.",
		"I organize my working time efficiently.",
		"I meet the deadlines I am given.",
		"I prioritize my tasks efficiently.",
		"I handle several tasks at the same time.",
		"I am rigorous in my work.",
		"I organize myself without constant supervision.",
		"I adapt quickly to a new environment.",
		"I accept constructive criticism.",
		"I stay calm in stressful situations.",
		"I work well under pressure.",
		"I adapt to organizational changes.",
		"I handle the unexpected without panicking.",
		"I am ready to step out of my comfort zone.",
	}

	questions := make([]models.Question, len(contents))
	for i, c := range contents {
		questions[i] = models.Question{Content: c, Options: yesOrNo, Weight: 1}
	}

	if err := g.Create(&questions).Error; err != nil {
		return err
	}
	log.Printf("seed: %d questions added", len(questions))
	return nil
}

func SeedUsers(g *gorm.DB) error {
	var count int64
	if err := g.Model(&models.User{}).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	hash, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	type spec struct {
		name, city, sector, role string
		skills                   []string
		year                     int
	}
	specs := []spec{
		{"Alice Martin", "Paris", "Communication", "candidate", []string{"Rédaction web", "Réseaux sociaux"}, 1998},
		{"Lucas Bernard", "Lyon", "Développement", "candidate", []string{"Go", "PostgreSQL", "Docker"}, 1996},
		{"Emma Petit", "Marseille", "Design", "candidate", []string{"Figma", "UX", "Illustration"}, 1999},
		{"Hugo Robert", "Toulouse", "Data", "candidate", []string{"Python", "SQL", "Pandas"}, 1997},
		{"Chloe Richard", "Nantes", "Marketing", "candidate", []string{"SEO", "Google Ads"}, 2000},
		{"Nathan Durand", "Nice", "Développement", "candidate", []string{"Vue.js", "TypeScript"}, 1995},
		{"Lea Moreau", "Strasbourg", "RH", "candidate", []string{"Recrutement", "Paie"}, 1994},
		{"Tom Laurent", "Montpellier", "Support", "candidate", []string{"Zendesk", "Relation client"}, 2001},
		{"Manon Simon", "Bordeaux", "Comptabilité", "candidate", []string{"Sage", "Excel"}, 1993},
		{"Enzo Michel", "Lille", "Logistique", "candidate", []string{"SAP", "Gestion de stock"}, 1998},
		{"Camille Garcia", "Rennes", "Communication", "candidate", []string{"Événementiel", "Relations presse"}, 1999},
		{"Louis David", "Reims", "Développement", "candidate", []string{"Java", "Spring"}, 1996},
		{"Sarah Bertrand", "Le Havre", "Design", "candidate", []string{"Photoshop", "Branding"}, 1997},
		{"Jules Roux", "Dijon", "Data", "candidate", []string{"Power BI", "SQL"}, 2000},
		{"Ines Vincent", "Grenoble", "Marketing", "candidate", []string{"Content", "Emailing"}, 1995},
		{"Adam Fournier", "Angers", "Développement", "candidate", []string{"Node.js", "React"}, 1998},
		{"Julie Girard", "Nimes", "RH", "candidate", []string{"Formation", "SIRH"}, 1994},
		{"Raphael Bonnet", "Metz", "Support", "candidate", []string{"ITIL", "Diagnostic"}, 2001},
		{"Claire Dupont", "Tours", "Recrutement", "recruiter", nil, 1990},
		{"Marc Leroy", "Paris", "Direction", "admin", nil, 1985},
	}

	users := make([]models.User, len(specs))
	for i, s := range specs {
		email := strings.ToLower(strings.ReplaceAll(s.name, " ", ".")) + "@competences.fr"
		users[i] = models.User{
			Name:         s.name,
			Email:        email,
			PasswordHash: string(hash),
			BirthDay:     time.Date(s.year, 6, 15, 0, 0, 0, 0, time.UTC),
			Status:       string(models.StatusAdult),
			Skills:       pq.StringArray(s.skills),
			Sector:       s.sector,
			City:         s.city,
			Role:         s.role,
		}
	}

	if err := g.Create(&users).Error; err != nil {
		return err
	}
	log.Printf("seed: %d users added (mot de passe: password123)", len(users))
	return nil
}

func SeedAdmins(g *gorm.DB) error {
	hash, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	names := []string{"maryam", "kevser", "leo"}
	added := 0
	for _, name := range names {
		email := name + "@competences.fr"

		var count int64
		if err := g.Model(&models.User{}).Where("email = ?", email).Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			continue
		}

		admin := models.User{
			Name:         name,
			Email:        email,
			PasswordHash: string(hash),
			BirthDay:     time.Date(1995, 1, 1, 0, 0, 0, 0, time.UTC),
			Status:       string(models.StatusAdult),
			Sector:       "Direction",
			City:         "Paris",
			Role:         string(models.Admin),
		}
		if err := g.Create(&admin).Error; err != nil {
			return err
		}
		added++
	}

	if added > 0 {
		log.Printf("seed: %d admins added (mot de passe: password123)", added)
	}
	return nil
}
