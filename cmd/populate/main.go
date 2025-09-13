package main

import (
	"log"

	"andrei-api/config"
	"andrei-api/models"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(".env"); err != nil {
		log.Println("⚠️ No .env file found")
	}

	// Connect to database
	config.ConnectDatabase()

	log.Println("🌱 Iniciando población de la base de datos...")

	// Create demons
	demons := createDemons()
	log.Printf("✅ Creados %d demonios", len(demons))

	// Create network admins (victims)
	networkAdmins := createNetworkAdmins()
	log.Printf("✅ Creados %d administradores de red", len(networkAdmins))

	// Create reports
	reports := createReports(demons, networkAdmins)
	log.Printf("✅ Creados %d reportes", len(reports))

	// Create rewards and punishments
	rewards := createRewards(demons)
	log.Printf("✅ Creadas %d recompensas/castigos", len(rewards))

	// Create posts
	posts := createPosts(demons, networkAdmins)
	log.Printf("✅ Creados %d posts", len(posts))

	log.Println("🎉 Base de datos poblada exitosamente!")
	printSummary()
}

func createDemons() []models.User {
	demons := []models.User{
		{Username: "ShadowMaster", Email: "shadow@evil.com", Role: models.RoleDemon},
		{Username: "DarkLord666", Email: "darklord@evil.com", Role: models.RoleDemon},
		{Username: "ChaosReaper", Email: "chaos@evil.com", Role: models.RoleDemon},
		{Username: "VoidWhisperer", Email: "void@evil.com", Role: models.RoleDemon},
		{Username: "NightmareKing", Email: "nightmare@evil.com", Role: models.RoleDemon},
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("demon123"), bcrypt.DefaultCost)

	for i := range demons {
		demons[i].Password = string(hashedPassword)
		if err := config.DB.Create(&demons[i]).Error; err != nil {
			log.Printf("⚠️ Error creando demonio %s: %v", demons[i].Username, err)
		}
	}

	return demons
}

func createNetworkAdmins() []models.User {
	networkAdmins := []models.User{
		{Username: "AdminJohn", Email: "john.admin@company.com", Role: models.RoleNetworkAdmin},
		{Username: "TechSarah", Email: "sarah.tech@corp.com", Role: models.RoleNetworkAdmin},
		{Username: "NetMike", Email: "mike.net@enterprise.com", Role: models.RoleNetworkAdmin},
		{Username: "SysLisa", Email: "lisa.sys@organization.com", Role: models.RoleNetworkAdmin},
		{Username: "AdminCarlos", Email: "carlos.admin@business.com", Role: models.RoleNetworkAdmin},
		{Username: "TechAna", Email: "ana.tech@company.com", Role: models.RoleNetworkAdmin},
		{Username: "NetDavid", Email: "david.net@corporation.com", Role: models.RoleNetworkAdmin},
		{Username: "SysEmily", Email: "emily.sys@enterprise.com", Role: models.RoleNetworkAdmin},
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)

	for i := range networkAdmins {
		networkAdmins[i].Password = string(hashedPassword)
		if err := config.DB.Create(&networkAdmins[i]).Error; err != nil {
			log.Printf("⚠️ Error creando admin %s: %v", networkAdmins[i].Username, err)
		}
	}

	return networkAdmins
}

func createReports(demons, networkAdmins []models.User) []models.Report {
	reportData := []struct {
		demonIdx    int
		victimIdx   int
		title       string
		description string
		status      string
	}{
		{0, 0, "Initial Contact", "Successfully approached target AdminJohn during lunch break", "completed"},
		{0, 1, "Hypnosis Progress", "TechSarah showing signs of susceptibility to mind control", "in_progress"},
		{1, 2, "Infiltration Report", "NetMike has been compromised. Access to server room obtained", "completed"},
		{1, 3, "Resistance Detected", "SysLisa is fighting back. Need backup", "pending"},
		{2, 4, "Complete Control", "AdminCarlos is now fully under our influence", "completed"},
		{2, 5, "Social Engineering", "TechAna fell for the fake IT support call", "completed"},
		{3, 6, "Physical Access", "Gained entry to NetDavid's workstation", "in_progress"},
		{4, 7, "Email Compromise", "SysEmily's email account has been accessed", "completed"},
		{0, 2, "Follow-up Report", "NetMike requires continued monitoring", "in_progress"},
		{1, 0, "Status Update", "AdminJohn showing signs of rebellion", "pending"},
	}

	var reports []models.Report
	for _, data := range reportData {
		if data.demonIdx < len(demons) && data.victimIdx < len(networkAdmins) {
			report := models.Report{
				DemonID:     demons[data.demonIdx].ID,
				VictimID:    networkAdmins[data.victimIdx].ID,
				Title:       data.title,
				Description: data.description,
				Status:      data.status,
			}
			if err := config.DB.Create(&report).Error; err != nil {
				log.Printf("⚠️ Error creando reporte: %v", err)
			} else {
				reports = append(reports, report)
			}
		}
	}

	return reports
}

func createRewards(demons []models.User) []models.Reward {
	rewardData := []struct {
		demonIdx    int
		rewardType  models.RewardType
		title       string
		description string
		points      int
	}{
		{0, models.RewardTypeReward, "Excellent Infiltration", "Successfully compromised 3 network admins this week", 150},
		{1, models.RewardTypeReward, "Social Engineering Master", "Perfect execution of phishing campaign", 100},
		{2, models.RewardTypePunishment, "Sloppy Work", "Target escaped due to careless approach", -50},
		{3, models.RewardTypeReward, "Stealth Operation", "Infiltrated high-security facility undetected", 200},
		{4, models.RewardTypeReward, "Team Player", "Assisted other demons in their missions", 75},
		{1, models.RewardTypePunishment, "Blown Cover", "Identity revealed to security team", -100},
		{0, models.RewardTypeReward, "Innovation Award", "Developed new hypnosis technique", 125},
		{3, models.RewardTypePunishment, "Mission Failure", "Target successfully alerted authorities", -75},
		{2, models.RewardTypeReward, "Loyalty Bonus", "Unwavering dedication to the cause", 50},
		{4, models.RewardTypeReward, "Rapid Progress", "Completed mission ahead of schedule", 90},
	}

	var rewards []models.Reward
	for _, data := range rewardData {
		if data.demonIdx < len(demons) {
			reward := models.Reward{
				DemonID:     demons[data.demonIdx].ID,
				Type:        data.rewardType,
				Title:       data.title,
				Description: data.description,
				Points:      data.points,
			}
			if err := config.DB.Create(&reward).Error; err != nil {
				log.Printf("⚠️ Error creando recompensa: %v", err)
			} else {
				rewards = append(rewards, reward)
			}
		}
	}

	return rewards
}

func createPosts(demons, networkAdmins []models.User) []models.Post {
	// Get Andrei user
	var andrei models.User
	config.DB.Where("role = ?", models.RoleAndrei).First(&andrei)

	posts := []struct {
		authorID  *uint
		title     string
		body      string
		media     string
		anonymous bool
	}{
		// Andrei's posts
		{&andrei.ID, "Welcome to the New Order", "My loyal demons, the time has come to expand our dominion over the digital realm. Every network administrator captured brings us closer to total control!", "", false},
		{&andrei.ID, "Survival Tips for the Conquered", "Remember, resistance is futile. The sooner you accept your new reality, the easier it will be.", "", false},

		// Demon posts
		{&demons[0].ID, "Infiltration Techniques", "Brothers, I've discovered that posing as IT support is incredibly effective. Humans trust anyone who claims they can fix their computer problems.", "", false},
		{&demons[1].ID, "Social Engineering 101", "The key is patience. Build trust first, then strike when they least expect it.", "", false},
		{&demons[2].ID, "Meme: When You Successfully Hypnotize Your First Admin", "That feeling when they start questioning their own reality 😈", "https://example.com/evil-meme.jpg", false},
		{&demons[3].ID, "Advanced Mind Control", "The new hypnosis protocols are working wonderfully. Remember to practice your incantations daily.", "", false},
		{&demons[4].ID, "Network Vulnerabilities", "I've compiled a list of common security weaknesses. Use them wisely.", "", false},

		// Anonymous posts from Network Admins
		{nil, "RESIST THE DARKNESS", "Fellow administrators, do not give in to their influence! We must fight back against these digital demons!", "", true},
		{nil, "Underground Communication", "If you're reading this, know that there are still those of us who resist. Stay strong!", "", true},
		{nil, "Security Tips", "Change your passwords frequently and trust no one claiming to be IT support without proper verification!", "", true},
		{nil, "They Are Among Us", "I've seen them in the server rooms, whispering their dark incantations. Be vigilant!", "", true},
		{nil, "Meme: Network Admins vs Demons", "When you successfully block another infiltration attempt 💪", "https://example.com/resistance-meme.jpg", true},
	}

	var createdPosts []models.Post
	for _, postData := range posts {
		post := models.Post{
			AuthorID:  postData.authorID,
			Title:     postData.title,
			Body:      postData.body,
			Media:     postData.media,
			Anonymous: postData.anonymous,
		}
		if err := config.DB.Create(&post).Error; err != nil {
			log.Printf("⚠️ Error creando post: %v", err)
		} else {
			createdPosts = append(createdPosts, post)
		}
	}

	return createdPosts
}

func printSummary() {
	var stats models.PlatformStats
	config.DB.Model(&models.User{}).Count(&stats.TotalUsers)
	config.DB.Model(&models.User{}).Where("role = ?", models.RoleDemon).Count(&stats.TotalDemons)
	config.DB.Model(&models.User{}).Where("role = ?", models.RoleNetworkAdmin).Count(&stats.TotalNetworkAdmins)
	config.DB.Model(&models.Post{}).Count(&stats.TotalPosts)
	config.DB.Model(&models.Report{}).Count(&stats.TotalReports)

	log.Println("\n📊 ESTADÍSTICAS FINALES:")
	log.Printf("   👥 Total usuarios: %d", stats.TotalUsers)
	log.Printf("   👹 Demonios: %d", stats.TotalDemons)
	log.Printf("   👨‍💻 Network Admins: %d", stats.TotalNetworkAdmins)
	log.Printf("   📝 Posts: %d", stats.TotalPosts)
	log.Printf("   📊 Reportes: %d", stats.TotalReports)

	log.Println("\n🔑 CREDENCIALES DE PRUEBA:")
	log.Println("   👑 Andrei (Admin): andrei@evil.com / password123")
	log.Println("   👹 Demonio: shadow@evil.com / demon123")
	log.Println("   👨‍💻 Network Admin: john.admin@company.com / admin123")
}
