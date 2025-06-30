using System.ComponentModel.DataAnnotations;

namespace SimpleTodoAPI.Models
{
    public class Todo
    {
        public int Id { get; set; }
        
        [Required]
        [StringLength(200)]
        public string Title { get; set; } = string.Empty;
        
        public bool Completed { get; set; } = false;
        
        public DateTime CreatedAt { get; set; } = DateTime.UtcNow;
    }

    public class CreateTodoRequest
    {
        [Required]
        public string Title { get; set; } = string.Empty;
    }

    public class UpdateTodoRequest
    {
        public string? Title { get; set; }
        public bool? Completed { get; set; }
    }
}